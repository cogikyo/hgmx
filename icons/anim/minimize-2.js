import { animate, τmax } from "../motion.js";

export default function (el) {
  const arrows = el.querySelector("g");
  animate(arrows, { scale: [1.2, 0.9, 1] }, { duration: τmax });
}
