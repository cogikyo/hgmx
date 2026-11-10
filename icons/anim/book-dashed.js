import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const outline = ["M4 19.5V14", "M4 10V8.5", "M4 4.5A", "M12 2h", "M17.5 2H", "M20 8.5V", "M20 14v", "M17.5 22H", "M12 22h", "M8 22H", "M12 17h"].map((d) =>
    el.querySelector(`path[d^="${d}"]`),
  );
  animate(outline, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });
  animate(outline, { translateY: [0, -3, 1, 0] }, { duration: τmax });
}
