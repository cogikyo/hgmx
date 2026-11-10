import { animate, τmax } from "../motion.js";

export default function (el) {
  const tie = el.querySelector("g");
  animate(tie, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
