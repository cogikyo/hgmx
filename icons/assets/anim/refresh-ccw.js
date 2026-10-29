import { animate, τmax } from "../motion.js";

export default function (el) {
  const arrows = el.querySelector("g");
  animate(arrows, { rotate: [0, -360], scale: [1, 0.85, 1] }, { duration: τmax });
}
