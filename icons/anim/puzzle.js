import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const piece = el.querySelector("path");
  animate(piece, { translateY: [-4, 1, 0], rotate: [-15, 5, 0], opacity: [0, 1] }, { duration: τmax });
  animate(piece, { scale: [1, 1.1, 1] }, { duration: τmin, delay: τmid });
}
