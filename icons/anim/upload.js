import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [6, -2, 0], opacity: [0, 1] }, { duration: τmax });

  const tray = el.querySelector('path[d^="M21 15"]');
  animate(tray, { scaleY: [1, 0.7, 1] }, { duration: τmin });
}
