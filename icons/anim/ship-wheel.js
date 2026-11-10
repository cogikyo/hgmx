import { animate, τmax } from "../motion.js";

export default function (el) {
  const helm = el.querySelector("g");
  animate(helm, { rotate: [0, -30, 20, -8, 0] }, { duration: τmax });
}
