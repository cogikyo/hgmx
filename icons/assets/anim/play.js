import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const play = el.querySelector("path");
  animate(play, { translateX: [0, 4, -1, 0] }, { duration: τmax });
  animate(play, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
