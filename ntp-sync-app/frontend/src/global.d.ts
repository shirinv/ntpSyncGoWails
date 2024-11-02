// frontend/src/global.d.ts

export {};

declare global {
  interface Window {
    backend: any;
    runtime: any;
  }
}
