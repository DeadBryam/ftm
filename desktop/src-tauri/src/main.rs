// Prevents additional console window on Windows in release
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use log::info;

fn main() {
    env_logger::Builder::from_env(env_logger::Env::default().default_filter_or("info")).init();
    info!("Starting Foundry Tunnel Manager Desktop");

    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .setup(|app| {
            ftm_desktop_lib::setup_app(app)?;
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            ftm_desktop_lib::commands::get_server_url,
            ftm_desktop_lib::commands::get_web_port
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
