import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [0, -4, 1, 0] }, { duration: τmax });

  const tray = el.querySelector(":scope > path");
  animate(tray, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
