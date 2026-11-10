import { animate, τmax } from "../motion.js";

export default function (el) {
  const rat = el.querySelector("g");
  animate(rat, { rotate: [0, 4, 0, 3, 0] }, { duration: τmax });
}
