import { animate, stagger } from "../motion.js";

export default function (el) {
  const waves = ["M8.5", "M5", "M12"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
