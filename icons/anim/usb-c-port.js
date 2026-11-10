import { animate, τmax } from "../motion.js";

export default function (el) {
  const slot = el.querySelector("path");
  animate(slot, { scaleX: [0, 1.25, 1] }, { duration: τmax });
}
