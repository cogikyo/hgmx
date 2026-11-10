import { animate, τmax } from "../motion.js";

export default function (el) {
  const hook = el.querySelector("g");
  animate(hook, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
