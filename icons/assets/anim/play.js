import { animate, τmax } from "../motion.js";

export default function (el) {
  const play = el.querySelector("path");
  animate(play, { translateX: [0, 4, -1, 0] }, { duration: τmax });
}
