import { animate, τmax } from "../motion.js";

export default function (el) {
  const whisk = el.querySelector("g");
  animate(whisk, { rotate: [0, -12, 10, -6, 0] }, { duration: τmax });
}
