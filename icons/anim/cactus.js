import { animate, τmax } from "../motion.js";

export default function (el) {
  const cactus = el.querySelector("g");
  animate(cactus, { rotate: [0, -6, 4, -2, 0] }, { duration: τmax });
}
