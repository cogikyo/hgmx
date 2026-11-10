import { animate, τmax } from "../motion.js";

export default function (el) {
  const chair = el.querySelector("g");
  animate(chair, { rotate: [0, -8, 6, -3, 0] }, { duration: τmax });
}
