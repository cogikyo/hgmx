import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("path");
  animate(arrow, { translateX: [0, -4, 1, 0] }, { duration: τmax });

  const dash = el.querySelector('path[d="M20 9v6"]');
  animate(dash, { translateX: [0, 1, 0] }, { duration: τmin });
}
