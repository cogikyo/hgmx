import { animate, τmax } from "../motion.js";

export default function (el) {
  const muzzle = el.querySelector("g");
  animate(muzzle, { scaleY: [1, 1.2, 0.95, 1] }, { duration: τmax });
}
