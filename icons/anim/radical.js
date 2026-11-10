import { animate, τmax } from "../motion.js";

export default function (el) {
  const radical = el.querySelector("path");
  animate(radical, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
