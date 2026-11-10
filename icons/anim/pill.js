import { animate, τmax } from "../motion.js";

export default function (el) {
  const pill = el.querySelector("g");
  animate(pill, { translateY: [-5, 1, 0], rotate: [-90, 10, 0], opacity: [0, 1] }, { duration: τmax });
}
