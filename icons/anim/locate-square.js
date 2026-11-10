import { animate, τmax } from "../motion.js";

export default function (el) {
  const reticle = el.querySelector("g");
  animate(reticle, { scale: [1.5, 0.9, 1], opacity: [0, 1] }, { duration: τmax });
}
