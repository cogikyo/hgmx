import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M5.254"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const second = el.querySelector('path[d^="M22 22"]');
  animate(second, { pathLength: [0, 1], opacity: [0, 1], translateX: [3, 0] }, { duration: τmid, delay: τmin });

  const body = el.querySelector('path[d^="M17 22"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
