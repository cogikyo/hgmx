import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("path");
  animate(arrow, { translateY: [0, -4, 1, 0] }, { duration: τmax });

  const dash = el.querySelector('path[d="M9 20h6"]');
  animate(dash, { translateY: [0, 1, 0] }, { duration: τmin });
}
