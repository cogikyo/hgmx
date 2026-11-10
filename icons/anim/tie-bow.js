import { animate, τmax } from "../motion.js";

export default function (el) {
  const bow = el.querySelector("g");
  animate(bow, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
