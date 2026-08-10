import { mount, unmount } from 'svelte';
import MiniPanel from '$lib/components/MiniPanel.svelte';

const PIP_WIDTH = 340;
const PIP_HEIGHT = 220;

interface DocumentPictureInPicture {
  requestWindow(options?: { width?: number; height?: number }): Promise<Window>;
  window: Window | null;
}

function documentPip(): DocumentPictureInPicture | null {
  if (typeof window === 'undefined') return null;
  return (window as unknown as { documentPictureInPicture?: DocumentPictureInPicture })
    .documentPictureInPicture ?? null;
}

export function supportsDocumentPip(): boolean {
  return documentPip() !== null;
}

function copyStyles(target: Window) {
  for (const sheet of Array.from(document.styleSheets)) {
    try {
      const css = Array.from(sheet.cssRules)
        .map((rule) => rule.cssText)
        .join('');
      const style = target.document.createElement('style');
      style.textContent = css;
      target.document.head.appendChild(style);
    } catch {
      if (!sheet.href) continue;
      const link = target.document.createElement('link');
      link.rel = 'stylesheet';
      link.href = sheet.href;
      target.document.head.appendChild(link);
    }
  }
}

function copyTheme(target: Window) {
  const source = document.documentElement;
  for (const attribute of Array.from(source.attributes)) {
    target.document.documentElement.setAttribute(attribute.name, attribute.value);
  }
}

export async function openDocumentPip(tunnelId: string): Promise<void> {
  const pip = documentPip();
  if (!pip) throw new Error('Document Picture-in-Picture is unavailable');

  const target = await pip.requestWindow({ width: PIP_WIDTH, height: PIP_HEIGHT });

  copyStyles(target);
  copyTheme(target);
  target.document.body.style.margin = '0';

  const panel = mount(MiniPanel, { target: target.document.body, props: { tunnelId } });

  target.addEventListener('pagehide', () => void unmount(panel), { once: true });
}
