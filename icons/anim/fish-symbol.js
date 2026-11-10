import { animate, τmax } from "../motion.js";

export default function (el) {
  const fish = el.querySelector("path");
  animate(fish, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
