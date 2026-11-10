import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const brick = el.querySelector("g");
  animate(brick, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: τmax });
  animate(brick, { scaleY: [1, 0.9, 1] }, { duration: τmin, delay: τmid });
}
