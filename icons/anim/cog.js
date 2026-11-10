import { animate, τmax } from "../motion.js";

export default function (el) {
  const wheel = el.querySelector("g");
  animate(wheel, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmax });
}
