import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const stick = el.querySelector("g");
  animate(stick, { rotate: [0, -15, 10, 0] }, { duration: τmax });
  const button = el.querySelector('path[d="M6 15v-2"]');
  animate(button, { scaleY: [1, 0.5, 1] }, { duration: τmid, delay: τmin });
}
