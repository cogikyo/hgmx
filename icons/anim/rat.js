import { animate, τmax } from "../motion.js";

export default function (el) {
  const rat = el.querySelector("g");
  animate(rat, { rotate: [0, 8, 0, 6, 0] }, { duration: τmax });
}
