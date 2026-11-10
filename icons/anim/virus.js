import { animate, τmax } from "../motion.js";

export default function (el) {
  const virus = el.querySelector("g");
  animate(virus, { rotate: [0, 90], scale: [1, 0.9, 1] }, { duration: τmax });
}
