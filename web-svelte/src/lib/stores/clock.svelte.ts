let now = $state(Date.now());
let timer: ReturnType<typeof setInterval> | null = null;
let subscribers = 0;

export function useClock() {
  return {
    get now() {
      return now;
    },

    subscribe() {
      subscribers += 1;
      if (!timer) {
        timer = setInterval(() => (now = Date.now()), 1000);
      }
      return () => {
        subscribers -= 1;
        if (subscribers <= 0 && timer) {
          clearInterval(timer);
          timer = null;
          subscribers = 0;
        }
      };
    },
  };
}
