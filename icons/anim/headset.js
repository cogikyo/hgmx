import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const headset = el.querySelector("g");
  animate(headset, { scaleX: [1, 0.8, 1.1, 1] }, { duration: τmax });

  const boom = el.querySelector('path[d^="M21 16"]');
  animate(boom, { rotate: [0, -10, 0] }, { duration: τmid, delay: τmin });
}
