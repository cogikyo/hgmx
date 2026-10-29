import { animate, τmax } from "../motion.js";

export default function (el) {
  const flake = el.querySelector("g");
  animate(flake, { translateY: [-4, 1, 0], rotate: [0, 180], opacity: [0, 1] }, { duration: τmax });
}
