import { animate, τmax } from "../motion.js";

export default function (el) {
  const squirrel = el.querySelector("g");
  animate(squirrel, { rotate: [0, -8, 0, -6, 0] }, { duration: τmax });
}
