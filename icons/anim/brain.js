import { animate, τmax } from "../motion.js";

export default function (el) {
  const brain = el.querySelector("g");
  animate(brain, { scale: [1, 0.9, 1.1, 1] }, { duration: τmax });
}
