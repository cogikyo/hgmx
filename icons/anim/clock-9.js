import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const face = el.querySelector("circle");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });
  const hands = el.querySelector("path");
  animate(hands, { rotate: [0, 360] }, { duration: τmax });
}
