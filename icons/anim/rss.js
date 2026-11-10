import { animate, stagger } from "../motion.js";

export default function (el) {
  const elements = Array.from(el.querySelectorAll("circle, path"));
  animate(
    elements,
    { scale: [0, 1.5, 1], opacity: [0, 1] },
    { duration: 0.3, delay: stagger(0.09) },
  );
}
