import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pen = el.querySelector('path[d*="21.378"]');
  animate(pen, { rotate: [-15, 15, -10, 0], translateX: [0, -5, 0] }, { duration: τmax });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, 1, -2, 0], opacity: [1] }, { duration: τmin, delay: τmid });

  const body = el.querySelector('path[d*="15H7"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
