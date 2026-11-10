import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const car = el.querySelectorAll("path, circle");
  animate(car, { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const body = el.querySelector('path[d^="M19 17"]');
  animate(body, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
