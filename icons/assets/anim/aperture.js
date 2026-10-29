import { animate, τmax } from "../motion.js";

export default function (el) {
  const blades = el.querySelector("g");
  animate(blades, { rotate: [0, 360], scale: [1, 0.6, 1] }, { duration: τmax });
}
