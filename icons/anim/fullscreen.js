import { animate, τmax } from "../motion.js";

export default function (el) {
  const screen = el.querySelector("rect");
  animate(screen, { scale: [0.7, 1.2, 1] }, { duration: τmax });
}
