import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const finders = el.querySelectorAll("rect");
  animate(finders, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bits = ["M12 3h", "M12 7v", "M3 12h", "M16 12h", "M21 12v", "M21 16h", "M12 16v", "M12 21v", "M21 21v"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(bits, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
