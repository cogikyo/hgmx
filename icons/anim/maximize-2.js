import { animate, τmax } from "../motion.js";

export default function (el) {
  const arrows = el.querySelector("g");
  animate(arrows, { scale: [0.8, 1.1, 1] }, { duration: τmax });
}
