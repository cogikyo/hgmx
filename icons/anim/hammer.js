import { animate, τmax } from "../motion.js";

export default function (el) {
  const hammer = el.querySelector("g");
  animate(hammer, { rotate: [0, -25, 6, 0] }, { duration: τmax });
}
