import { animate, τmax } from "../motion.js";

export default function (el) {
  const bunch = el.querySelector("g");
  animate(bunch, { scale: [1, 0.9, 1.05, 1] }, { duration: τmax });
}
