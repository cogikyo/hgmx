import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const car = el.querySelector("g");
  animate(car, { scale: [0.5, 1.1, 1], opacity: [0, 1] }, { duration: τmax });

  const body = el.querySelectorAll('path[d^="m21 8"], rect, path[d$="h.01"]');
  animate(body, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
