use std::path::PathBuf;
use std::process::Stdio;
use std::time::Duration;
use tauri::{Manager, WebviewUrl};
use tokio::time::sleep;

fn main() {
    tauri::Builder::default()
        .setup(|app| {
            let app_handle = app.handle().clone();
            
            tauri::async_runtime::spawn(async move {
                start_go_backend(&app_handle).await;
                
                for _ in 0..50 {
                    sleep(Duration::from_millis(100)).await;
                    if check_server().await {
                        break;
                    }
                }
                
                let window = tauri::WebviewWindowBuilder::new(
                    &app_handle,
                    "main",
                    WebviewUrl::External("http://localhost:8080".parse().unwrap())
                )
                .title("Foundry Tunnel Manager")
                .inner_size(1200.0, 800.0)
                .min_inner_size(900.0, 600.0)
                .center()
                .build()
                .unwrap();
            });
            
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error running app");
}

async fn start_go_backend(app_handle: &tauri::AppHandle) {
    let resource_path = app_handle.path().resource_dir().unwrap();
    
    #[cfg(target_os = "windows")]
    let binary_name = "go-backend.exe";
    #[cfg(not(target_os = "windows"))]
    let binary_name = "go-backend";
    
    let binary_path = if cfg!(dev) {
        PathBuf::from("../go-api/cmd/server").join(binary_name)
    } else {
        resource_path.join(binary_name)
    };
    
    let _child = tokio::process::Command::new(&binary_path)
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn();
}

async fn check_server() -> bool {
    match tokio::net::TcpStream::connect("127.0.0.1:8080").await {
        Ok(_) => true,
        Err(_) => false,
    }
}
