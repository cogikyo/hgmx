import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const image = el.querySelector("g[transform-origin]");
  animate(image, { scale: [0.75, 1.1, 1] }, { duration: τmax });

  const arrow = el.querySelector("g:not([transform-origin])");
  animate(arrow, { translateX: [0, 2, -1, 0], translateY: [0, -2, 1, 0] }, { duration: τmid, delay: τmin });
}
