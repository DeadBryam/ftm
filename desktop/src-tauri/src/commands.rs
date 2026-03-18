use std::env;

const WEB_PORT_ENV: &str = "FOUNDRY_TUNNEL_WEB_PORT";

#[tauri::command]
pub fn get_server_url() -> String {
    let port = env::var(WEB_PORT_ENV).unwrap_or_else(|_| "8080".to_string());
    format!("http://localhost:{}/", port)
}

#[tauri::command]
pub fn get_web_port() -> u16 {
    env::var(WEB_PORT_ENV)
        .ok()
        .and_then(|p| p.parse().ok())
        .unwrap_or(8080)
}
