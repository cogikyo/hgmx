import { animate, stagger } from "../motion.js";

export default function (el) {
  const waves = [...el.querySelectorAll("path")];
  animate(waves.reverse(), { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
