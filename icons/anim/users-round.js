import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const second = el.querySelector('path[d^="M22 20"]');
  animate(second, { pathLength: [0, 1], opacity: [0, 1], translateX: [3, 0] }, { duration: τmid });

  const body = el.querySelector('path[d^="M18 21"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
