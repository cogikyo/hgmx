import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [6, -2, 0], opacity: [0, 1] }, { duration: τmax });

  const line = el.querySelector('path[d="M21 19V5"]');
  animate(line, { translateX: [0, 1, 0] }, { duration: τmin });
}
