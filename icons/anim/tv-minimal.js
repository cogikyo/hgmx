import { animate, τmid } from "../motion.js";

export default function (el) {
  const screen = el.querySelector("rect");
  animate(screen, { scaleY: [0.1, 1.1, 1], opacity: [0.5, 1] }, { duration: τmid });
}
